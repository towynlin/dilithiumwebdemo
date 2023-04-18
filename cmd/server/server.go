package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudflare/circl/sign/eddilithium2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/towynlin/dilithiumwebdemo/eddilithium2jwt"
)

func customJWTAuth(conn *pgx.Conn) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			extractor := request.BearerExtractor{}
			tokenString, err := extractor.ExtractToken(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			algs := []string{eddilithium2jwt.SigningMethodEdDilithium2.Alg()}
			parser := jwt.NewParser(jwt.WithIssuedAt(), jwt.WithValidMethods(algs))
			_, err = parser.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if err := eddilithium2jwt.ValidateIssuedAt(t); err != nil {
					return nil, err
				}
				issuer, err := t.Claims.GetIssuer()
				if err != nil {
					return nil, err
				}
				userID, err := uuid.Parse(issuer)
				if err != nil {
					return nil, err
				}
				row := conn.QueryRow(context.Background(), "select pubkey from users where id=$1", userID.String())
				pubKeyBytes := make([]byte, eddilithium2.PublicKeySize)
				if err = row.Scan(&pubKeyBytes); err != nil {
					return nil, err
				}
				var pubKey eddilithium2.PublicKey
				if err = pubKey.UnmarshalBinary(pubKeyBytes); err != nil {
					return nil, err
				}
				c.Set("user", userID)
				return pubKey, nil
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			return next(c)
		}
	}
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customJWTAuth(conn))

	e.GET("/jobs", getJobsHandler(conn))
	e.POST("/jobs", postJobHandler(conn))
	e.DELETE("/jobs/:id", deleteJobHandler(conn))

	e.Logger.Fatal(e.Start(":1323"))
}

type JobsResponse struct {
	Jobs []uuid.UUID `json:"jobs"`
}

type JobResponse struct {
	Id uuid.UUID `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func getJobsHandler(conn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		rows, _ := conn.Query(context.Background(), "select id from jobs where user_id=$1", user)
		ids, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (uuid.UUID, error) {
			var id uuid.UUID
			err := row.Scan(&id)
			return id, err
		})
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		}
		return c.JSON(http.StatusOK, JobsResponse{Jobs: ids})
	}
}

func postJobHandler(conn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := uuid.New()
		user := c.Get("user")
		_, err := conn.Exec(context.Background(), "insert into jobs (id, user_id) values ($1, $2)", id, user)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		}
		return c.JSON(http.StatusCreated, JobResponse{Id: id})
	}
}

func deleteJobHandler(conn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		maybeId := c.Param("id")
		id, err := uuid.Parse(maybeId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid job ID"})
		}
		user := c.Get("user")
		result, err := conn.Exec(context.Background(), "delete from jobs where id=$1 and user_id=$2", id, user)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		}
		if result.RowsAffected() == 0 {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Unknown job"})
		}
		return c.JSON(http.StatusNoContent, nil)
	}
}
