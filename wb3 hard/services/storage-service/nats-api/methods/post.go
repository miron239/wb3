package methods

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thoas/go-funk"
	"go-microservices/libs/cockroach"
	"go-microservices/libs/json_codec"
	"go-microservices/libs/nats"
	"go-microservices/services/client-service/rest"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetPosts(param map[string]interface{}, clientID string) {
	posts, err := cockroach.GetAllPosts()
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "GetPosts",
			ErrorCode: 10,
		})
		return
	}

	nats.Publish(nats.NatsMessage{
		Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
		ClientID:  clientID,
		IsSuccess: true,
		Details: map[string]interface{}{
			"posts": funk.Map(posts, func(post cockroach.Post) map[string]interface{} {
				return map[string]interface{}{
					"id":          strconv.FormatInt(post.ID, 10),
					"UserId":      post.UserId,
					"Spp":         post.Spp,
					"ShippingFee": post.ShippingFee,
					"ReturnFee":   post.ReturnFee,
					"date":        post.CreatedAt.String(),
				}
			}),
		},
		Method: "GetPosts",
	})
}

func GetPost(param map[string]interface{}, clientID string) {
	id, err := getID("DeletePost", param, clientID)
	if err != nil {
		return
	}

	post, err := cockroach.GetPost(id)
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "GetPost",
			ErrorCode: 10,
		})
		return
	}

	nats.Publish(nats.NatsMessage{
		Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
		ClientID:  clientID,
		IsSuccess: true,
		Details: map[string]interface{}{
			"id":          strconv.FormatInt(post.ID, 10),
			"UserId":      post.UserId,
			"Spp":         post.Spp,
			"ShippingFee": post.ShippingFee,
			"ReturnFee":   post.ReturnFee,
			"date":        post.CreatedAt.String(),
		},
		Method: "GetPost",
	})
}

func NewPost(param map[string]interface{}, clientID string) {
	userId, err := json_codec.GetString("UserId", param)
	postpaidLimit, err := json_codec.GetInt64("PostpaidLimit", param)
	spp, err := json_codec.GetInt64("Spp", param)
	shippingFee, err := json_codec.GetInt64("ShippingFee", param)
	returnFee, err := json_codec.GetInt64("ReturnFee", param)
	//date, err := json_codec.GetInt64("date", param)

	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "NewPost",
			ErrorCode: 10,
		})
		return
	}

	post, err := cockroach.CreatePost(cockroach.Post{
		UserId:        userId,
		PostpaidLimit: int(postpaidLimit),
		Spp:           int(spp),
		ShippingFee:   int(shippingFee),
		ReturnFee:     int(returnFee),
	})
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "NewPost",
			ErrorCode: 10,
		})
		return
	}

	nats.Publish(nats.NatsMessage{
		Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
		ClientID:  clientID,
		IsSuccess: true,
		Details: map[string]interface{}{
			"id":            strconv.FormatInt(post.ID, 10),
			"UserId":        post.UserId,
			"PostpaidLimit": post.PostpaidLimit,
			"Spp":           post.Spp,
			"ShippingFee":   post.ShippingFee,
			"ReturnFee":     post.ReturnFee,
			"date":          post.CreatedAt.String(),
		},
		Method: "NewPost",
	})
}

func UpdatePost(param map[string]interface{}, clientID string) {
	userId, err := json_codec.GetString("UserId", param)
	postpaidLimit, err := json_codec.GetInt64("PostpaidLimit", param)
	spp, err := json_codec.GetInt64("Spp", param)
	shippingFee, err := json_codec.GetInt64("ShippingFee", param)
	returnFee, err := json_codec.GetInt64("ReturnFee", param)
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "UpdatePost",
			ErrorCode: 10,
		})
		return
	}
	id, err := getID("DeletePost", param, clientID)
	if err != nil {
		return
	}

	err = cockroach.UpdatePost(cockroach.Post{ID: id, UserId: userId, PostpaidLimit: int(postpaidLimit), Spp: int(spp),
		ShippingFee: int(shippingFee), ReturnFee: int(returnFee)}, cockroach.Post{ID: id, UserId: userId, PostpaidLimit: int(postpaidLimit), Spp: int(spp),
		ShippingFee: int(shippingFee), ReturnFee: int(returnFee)})
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "UpdatePost",
			ErrorCode: 10,
		})
		return
	}

	nats.Publish(nats.NatsMessage{
		Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
		ClientID:  clientID,
		IsSuccess: true,
		Details: map[string]interface{}{
			"success": true,
		},
		Method: "UpdatePost",
	})
}

func DeletePost(param map[string]interface{}, clientID string) {
	id, err := getID("DeletePost", param, clientID)
	if err != nil {
		return
	}

	err = cockroach.DeletePost(id)
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    "DeletePost",
			ErrorCode: 10,
		})
		return
	}

	nats.Publish(nats.NatsMessage{
		Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
		ClientID:  clientID,
		IsSuccess: true,
		Details: map[string]interface{}{
			"success": true,
		},
		Method: "DeletePost",
	})
}

func getID(method string, param map[string]interface{}, clientID string) (int64, error) {
	i, err := json_codec.GetString("id", param)
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    method,
			ErrorCode: 18,
			Details: map[string]interface{}{
				"field": "id",
			},
		})
		return 0, err
	}
	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		nats.Publish(nats.NatsMessage{
			Services:  []string{nats.CLIENT_SREVICE, nats.CLI_SERVICE},
			ClientID:  clientID,
			IsSuccess: false,
			Method:    method,
			ErrorCode: 20,
			Details: map[string]interface{}{
				"field": "id",
				"type":  "int",
			},
		})
		return 0, err
	}
	return id, nil
}

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKENPASS")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user rest.Post

		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
