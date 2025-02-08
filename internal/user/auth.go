package user

// generateJWT generates a JWT that encodes an identity.
//func GenerateJWT(identity entity.UserAuthIdentity) (string, error) {
//	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":       identity.GetID(),
//		"email":    identity.GetEmail(),
//		"username": identity.GetUserName(),
//		"role":     identity.GetRole(),
//		"exp":      time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
//	}).SignedString([]byte(s.signingKey))
//}
