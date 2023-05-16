// api/user.thrift
namespace go user
namespace py user

struct User {
    1: string uid (api.form="uid");
    2: string email (api.form="email");
    3: string nickname (api.form="nickname");
}

struct SignupRequest {
    1: string email (api.form="email", api.vd="email($)", go.tag="example:\"me@mail.com\"");
    2: string password (api.form="password");
    3: string nickname (api.form="nickname");
}

struct LoginRequest {
    1: string email (api.form="email", api.vd="email($)", go.tag="example:\"me@mail.com\"");
    2: string password (api.form="password");
}

struct LoginResponse {
    1: string token;
    2: string expire;
}

struct LogoutRequest {
}

struct LogoutResponse {
}

struct GetUserRequest {
    1: string uid (api.path="uid");
}

struct GetUserResponse {
}

struct GetUserMeRequest {
}

struct GetUserMeResponse {
}


service UserService {
    User signup(1: SignupRequest req) (api.post="/users/signup")
    LoginResponse login(1: LoginRequest req) (api.post="/users/login")
    LogoutResponse logout() (api.post="/users/logout")
}
