PHONY: help

help:
	@ECHO +------------------------------------------------------+
	@ECHO ^|    Validation Repository Tests - Make Commands       ^|
	@ECHO +------------------------------------------------------+
	@ECHO   ^> make help                                          - Display this help information
	@ECHO   ^> make validation-test                               - Run all tests
	@ECHO   ^> make TestCreateTransactionObserver                 - Run TestCreateTransactionObserver test
	@ECHO   ^> make TestGetTransactionObserver                    - Run TestGetTransactionObserver test
	@ECHO   ^> make TestUpdateStatus                              - Run TestUpdateStatus test
	@ECHO   ^> make TestPrecheckoutTransaction                    - Run TestPrecheckoutTransaction test
	@ECHO   ^> make TestDeleteTransactionObserver                 - Run TestDeleteTransactionObserver test
validation-test:
	go test -v ../test/validation/repository -run 'TestValidationRepositoryTestSuite'

TestCreateTransactionObserver:
	go test -v ../test/validation/repository -run 'TestValidationRepositoryTestSuite/TestCreateTransactionObserver'

TestGetTransactionObserver:
	go test -v ../test/validation/repository -run 'TestValidationRepositoryTestSuite/TestGetTransactionObserver'

TestUpdateStatus:
	go test -v ../test/validation/repository -run 'TestValidationRepositoryTestSuite/TestUpdateStatus'

TestPrecheckoutTransaction:
	go test -v ../test/validation/repository -run 'TestValidationRepositoryTestSuite/TestPrecheckoutTransaction'

TestDeleteTransactionObserver:
	go test -v ../test/validation/repository -run 'TestValidationRepositoryTestSuite/TestDeleteTransactionObserver'

