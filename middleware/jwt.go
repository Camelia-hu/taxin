package interceptor

//
//func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) {
//	// 排除登录接口
//	if info.FullMethod == "/user.UserService/Login" {
//		return handler(ctx, req)
//	}
//
//	// 验证 Token
//	token := extractToken(ctx)
//	claims, err := verifyJWT(token)
//	if err != nil {
//		return nil, status.Error(codes.Unauthenticated, "invalid token")
//	}
//
//	// 注入用户信息到上下文
//	ctx = context.WithValue(ctx, "user_id", claims.UserID)
//	return handler(ctx, req)
//}
