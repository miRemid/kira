# File Service

The File Service need Another unique token for calling the api request. Every user would get his file's token when they signup the websites
and if he needs refresh, he can call the refresh api to refresh his token.

The File Service maintain two database tables
First is the file-meta infomation table, which is used to store the file's meta infomation about name, size, userid
Second is the token-user table, which stores the userid and token


### Functions

These are the apis which don't need token but only support rpc request
- GenerateToken(userid) Token
- RefreshToken(userid) Token 

This api would return the image's body data
- GetImage(fileID) bytes

- UploadFile(token, file's meta infomation, file's body data) FileModel

- DeleteFile(token, fileID)

# Auth Service

The Auth Service is the "JWT Token generate and valid" service for User Service.

### Functions
All of Functions only support RPC request

- GenerateToken(userid, role) JwtToken
- RefreshToken(tokenString) JwtToken
- Valid(tokenString) ValidResponse

# User Service

The User Service is the 

### Functions

- Signup(username, password) SignupResponse
- Signin(username, password) SigninResponse
- GetUserInfomation(token) UserModel
- ChangePassword(token, old, new) ChangeResponse

# Gateway Service

Gateway Service is the entrance of all services, using JWT for valding and casbin for role's control

### Functions
- JWT Validation
- Casbin Check

# Site Service

Site Service is the service for downloading image and check site's summary datas

- GetImage(fileID) bytes
- State() StateResponse