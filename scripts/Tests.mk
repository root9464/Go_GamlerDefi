.PHONY: help referral-test referral-test-TestGetPaymentOrdersByAuthorID referral-test-TestGetAllPaymentOrders referral-test-TestCreatePaymentOrder referral-test-TestGetPaymentOrdersByAuthorID_Empty

help:
	@ECHO +------------------------------------------------------+
	@ECHO ^|    Referral Repository Tests - Make Commands         ^|
	@ECHO +------------------------------------------------------+
	@ECHO   ^> make referral-test                                       - Run all tests for TestReferralRepositoryTestSuite
	@ECHO   ^> make referral-test-TestGetPaymentOrdersByAuthorID        - Run the specific test for TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID
	@ECHO   ^> make referral-test-TestGetAllPaymentOrders               - Run the specific test for TestReferralRepositoryTestSuite/TestGetAllPaymentOrders
	@ECHO   ^> make referral-test-TestCreatePaymentOrder                - Run the specific test for TestReferralRepositoryTestSuite/TestCreatePaymentOrder
	@ECHO   ^> make referral-test-TestGetPaymentOrdersByAuthorID_Empty  - Run the specific test for TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID_Empty
	@ECHO   ^> make help                                                - Display this help information

referral-test:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite'

referral-test-TestGetPaymentOrdersByAuthorID:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID'

referral-test-TestGetAllPaymentOrders:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetAllPaymentOrders'

referral-test-TestCreatePaymentOrder:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestCreatePaymentOrder'

referral-test-TestGetPaymentOrdersByAuthorID_Empty:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID_Empty'