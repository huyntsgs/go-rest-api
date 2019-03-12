
The reposity is about basic CRUD RESTful API with Gin-gonic framework and MySql as database store.

When working on RESTul api projects, we tend try to make them work and not much consider about how project adapts with the changes.

In this project, I organize source code structure in three main layers and explain how it can adapt with requirement change.

#### models 

Contains all defination of data structures and used by all other layers. There are also the validation functions for data structures. 
These validation functions are called in the api handles. As we can see the products, users model. 
If there is any other data structure such as orders, customers..., will be put on this layer.

#### store

Responsible for database related job such as query, insert, update or delete. No business logic is implemented here.
In store.go, I defined interfaces ProductStore, UserStore.

````go
type ProductStore interface {
	GetProducts(limit, offset int, lastId int64) ([]*models.Product, error)
	GetSingleProduct(productId int64) (*models.Product, error)
	DeleteProduct(productId int64) (int64, error)
	UpdateProduct(p *models.Product) (int64, error)
	CreateProduct(p *models.Product) (int64, error)
}

type UserStore interface {
	Register(user *models.User) error
	Login(user *models.User) (*models.User, error)
}
````

The productHandler uses ProductStore as input so it can make use the functions of struct that implements ProductStore. The same for userHandler.

The advantages of this are twofold. Firstly, it gives our code a really clean structure, but more importantly it also opens up the potential to mock our database for unit testing.

Here, in mysql_db.go we define struct MySqlDB which embedded with *sql.DB. MySqlDB implements all methods of ProductStore, UserStore. mysql_products.go contains all functions for interactive with database for products. 

If we want to change database to Postgres or MongoDB, how can we solve? Well, we can add more files such as mongdb_store.go, mongodb_products.go in this layer.
We will define MongoDB struct which embeds mgo.Session. In mongodb_products.go, MongoDB satifies the interface ProductStore also UserStore in mongodb_users.go as we have done with MySqlDB.

We obviously do not need to change any thing in the handles function and routers.

#### api

Contains all handles of api. 
authHandler is a middleware which parses and validates json web token JWT. authHandler can set claim information to context for later using in main handles.
productHandler contains all api for products. productHandler takes input is ProductStore inteface as we mentioned above.

#### router

In setup router part in main.go, I used gin and group common api url. Some apis need to be authenticated then will use the authentication middleware. 
The authorize uses crypto functions to encrypt/decrypt claim data as specified of JSON WEB TOKEN - JWT. The apis using authentication middleware are create, update and delete product.

````go
v1 := router.Group("/api/v1")
{
	// Get products with params limit offset likes /products?limit=50&lastId=0
	v1.GET("/products", productHandler.GetProducts())
	v1.GET("/products/:productId", productHandler.GetSingleProduct())
}
v1.POST("/users/register", userHandler.Register())
v1.POST("/users/login", userHandler.Login())
v1.Use(api.Authorize())
{
	v1.POST("/products", productHandler.CreateProduct())
	v1.PUT("/products", productHandler.UpdateProduct())
	v1.DELETE("/products/:productId", productHandler.DeleteProduct())
}
````

We have totally seven apis.

- GET /api/v1/products?limit=50&lastId=0
- GET /api/v1/products/productId
- POST /api/v1/users/register '{"userName":"huyntsgs", "email":"huyntsgs@gmail.com", "password":"pass123"}'
- POST /api/v1/users/login '{"email":"huyntsgs@gmail.com", "password":"pass123"}'

- POST /api/v1/products '{"productName":"Iphone Xs", "image":"iphonexs.jpg", "rate": 4}'
- PUT /api/v1/products '{"productId":1, productName":"Iphone Xss", "image":"iphonexs.jpg", "rate": 4}'
- DELETE /api/v1/products/productId


#### Start project

Make sure you have installed MySql database and setting database information correctly to .env file.

````bash
git clone https://github.com/huyntsgs/go-rest-api.git

chmod +x start.sh

cd go-rest-api/

./start.sh
````

In order to interactive with application, we can use curl or postman. For simply, I use curl.

````bash
curl -X POST -H "Content-Type:application/json" -d "{\"userName\":\"huyntsgs\", \"email\":\"huyntsgs@gmail.com\", \"password\":\"pass123\"}"  http://127.0.0.1:8081/api/v1/users/register

curl -X POST -H "Content-Type:application/json" -d "{\"email\":\"huyntsgs@gmail.com\", \"password\":\"pass123\"}"  http://127.0.0.1:8081/api/v1/users/login
````
In order to authenticate apis, we need to set token from login to Authorization header.

````bash
curl -X POST -H "Content-Type:application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imh1eW50c2dzQGdtYWlsLmNvbSIsImV4cCI6MTU1MjM3NzE3NSwidXNlcklkIjoxLCJ1c2VyTmFtZSI6Imh1eW50c2dzIn0.Prxp4FCa684f6wjXXL1jMuAciqXd8zLme7_lOhNhiwM" -d "{\"productName\":\"IphoneXs\", \"image\":\"iphonexs.jpg\", \"rate\": 4}"  http://127.0.0.1:8081/api/v1/products
````




