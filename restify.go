package restify

import (
	"fmt"

	"reflect"

	"net/http"

	"strconv"

	"strings"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
	"github.com/labstack/echo"
)

type (
	Restify struct {
		DB    *gorm.DB
		Group *echo.Group
	}
)

func (r *Restify) Register(interfaces ...interface{}) {
	for _, i := range interfaces {
		fmt.Println("Got", i)
		entity := reflect.TypeOf(i).Elem()
		entityName := strings.ToLower(inflection.Plural(entity.Name()))
		fmt.Printf("Registering %+v\n", entityName)
		entitySlice := reflect.SliceOf(entity)

		r.Group.GET("/"+entityName, func(c echo.Context) error {
			results := reflect.New(entitySlice).Interface()

			// TODO: Support filtering...
			if err := r.DB.Find(results).Error; err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			fmt.Printf("%+v\n", results)

			return c.JSON(http.StatusOK, results)
		})

		r.Group.POST("/"+entityName, func(c echo.Context) error {
			v := reflect.New(entity).Interface()
			err := c.Bind(v)

			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			err = r.DB.Save(v).Error
			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, v)
		})

		r.Group.PUT("/"+entityName+"/:id", func(c echo.Context) error {
			v := reflect.New(entity).Interface()
			err := c.Bind(v)

			id, err := strconv.ParseUint(c.Param("id"), 10, 64)

			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			reflect.ValueOf(v).Elem().FieldByName("ID").SetUint(id)

			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			err = r.DB.Save(v).Error
			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, v)
		})

		r.Group.PATCH("/"+entityName+"/:id", func(c echo.Context) error {
			v := reflect.New(entity).Interface()
			err := c.Bind(v)

			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			id, err := strconv.ParseUint(c.Param("id"), 10, 64)

			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			reflect.ValueOf(v).Elem().FieldByName("ID").SetUint(id)

			r.DB.Model(v).Updates(v)
			r.DB.Find(v, id)

			return c.JSON(http.StatusOK, v)
		})
	}
}

func New(g *echo.Group, db *gorm.DB) (r *Restify) {
	r = &Restify{
		Group: g,
		DB:    db,
	}
	return
}
