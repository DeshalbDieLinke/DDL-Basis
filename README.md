# DDL-Basis
Backend für DeshalbDieLinke.de

# Auth-flow
- Server generates a 10 Minute token for the email specified in the system env. Variable INIT_EMAIL with access level 0 (admin). This only happens once. 
- Using this token an admin account can be created> This account will be able to invite new users generating a token for their email + access Level. The function for this is located at utils.GenerateToken 
- Users will be able to register with email, password, and Token. The Email must match the one in the Token Claims. 
- Now users will be able to login, which will prompt the frontend to store the token in localstorage
- When requesting anything requiring auth the token will be supplied in the request header "Authorization". No Bearer tag. 
- When the token is invalid the user should be redirected to /login (frontend) 

## Features
- The user can check if logged in using the /auth/check api endpoint
- The user can login with just the token and email by leaving the password field empty. This happens automatically if a token is found !!Bug if a password is provided it will fail. 

# Archi

![Archi](./archi.png)
