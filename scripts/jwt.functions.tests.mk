.PHONY: help jwt-test TestGenerateKeyPair_Success TestGenerateKeyPair_InvalidUserData TestGenerateKeyPair_NilPrivateKey TestGenerateKeyPair_NilHelpers TestRefreshAccessToken_Success TestRefreshAccessToken_InvalidToken TestRefreshAccessToken_ExpiredToken TestRefreshAccessToken_MissingClaims

help:
	@ECHO +------------------------------------------------------+
	@ECHO |    JWT Functions Tests - Make Commands               |
	@ECHO +------------------------------------------------------+
	@ECHO   > make jwt-test                                     - Run all tests for JwtFuncsTestSuite
	@ECHO   > make TestGenerateKeyPair_Success                  - Run JwtFuncsTestSuite/TestGenerateKeyPair_Success
	@ECHO   > make TestGenerateAdminToken_Success               - Run JwtFuncsTestSuite/TestGenerateAdminToken_Success
	@ECHO   > make TestGenerateKeyPair_InvalidUserData          - Run JwtFuncsTestSuite/TestGenerateKeyPair_InvalidUserData
	@ECHO   > make TestGenerateKeyPair_NilPrivateKey            - Run JwtFuncsTestSuite/TestGenerateKeyPair_NilPrivateKey
	@ECHO   > make TestGenerateKeyPair_NilHelpers               - Run JwtFuncsTestSuite/TestGenerateKeyPair_NilHelpers
	@ECHO   > make TestRefreshAccessToken_Success               - Run JwtFuncsTestSuite/TestRefreshAccessToken_Success
	@ECHO   > make TestRefreshAccessToken_InvalidToken          - Run JwtFuncsTestSuite/TestRefreshAccessToken_InvalidToken
	@ECHO   > make TestRefreshAccessToken_ExpiredToken          - Run JwtFuncsTestSuite/TestRefreshAccessToken_ExpiredToken
	@ECHO   > make TestRefreshAccessToken_MissingClaims         - Run JwtFuncsTestSuite/TestRefreshAccessToken_MissingClaims
	@ECHO   > make help                                        - Display this help information

jwt-test:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite'

TestGenerateKeyPair_Success:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestGenerateKeyPair_Success'

TestGenerateAdminToken_Success:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestGenerateAdminToken_Success'

TestGenerateKeyPair_InvalidUserData:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestGenerateKeyPair_InvalidUserData'

TestGenerateKeyPair_NilPrivateKey:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestGenerateKeyPair_NilPrivateKey'

TestGenerateKeyPair_NilHelpers:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestGenerateKeyPair_NilHelpers'

TestRefreshAccessToken_Success:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestRefreshAccessToken_Success'

TestRefreshAccessToken_InvalidToken:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestRefreshAccessToken_InvalidToken'

TestRefreshAccessToken_ExpiredToken:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestRefreshAccessToken_ExpiredToken'

TestRefreshAccessToken_MissingClaims:
	go test -v ../test/jwt/functions -run 'TestJwtFuncsTestSuite/TestRefreshAccessToken_MissingClaims'