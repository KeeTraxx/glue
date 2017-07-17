package glue

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

type Tester struct {
	gorm.Model
	Name string `json:"name"`
}

// Glue glues an echo.Group together with a *gorm.DB
func Glue(g *echo.Group, db *gorm.DB, interfaces ...interface{}) error {
	for _, i := range interfaces {
		entity := reflect.TypeOf(i).Elem()
		entityName := strings.ToLower(inflection.Plural(entity.Name()))
		entitySlice := reflect.SliceOf(entity)

		fmt.Printf("Registering GET /%v \n", entityName)
		g.GET("/"+entityName, func(c echo.Context) error {
			results := reflect.New(entitySlice).Interface()

			// TODO: Support filtering...
			if err := db.Find(results).Error; err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			fmt.Printf("%+v\n", results)

			return c.JSON(http.StatusOK, results)
		})

		fmt.Printf("Registering POST /%v \n", entityName)
		g.POST("/"+entityName, func(c echo.Context) error {
			v := reflect.New(entity).Interface()
			err := c.Bind(v)

			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			err = db.Save(v).Error
			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, v)
		})

		fmt.Printf("Registering PUT /%v \n", entityName)
		g.PUT("/"+entityName+"/:id", func(c echo.Context) error {
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

			err = db.Save(v).Error
			if err != nil {
				fmt.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, v)
		})

		fmt.Printf("Registering PATCH /%v \n", entityName)
		g.PATCH("/"+entityName+"/:id", func(c echo.Context) error {
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

			db.Model(v).Updates(v)
			db.Find(v, id)

			return c.JSON(http.StatusOK, v)
		})
	}
	return nil
}
