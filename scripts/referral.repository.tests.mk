.PHONY: help referral-test TestGetPaymentOrdersByAuthorID TestGetAllPaymentOrders TestCreatePaymentOrder TestGetPaymentOrdersByAuthorID_Empty

help:
	@ECHO +------------------------------------------------------+
	@ECHO ^|    Referral Repository Tests - Make Commands         ^|
	@ECHO +------------------------------------------------------+
	@ECHO   ^> make referral-test                                 - Run all tests for TestReferralRepositoryTestSuite
	@ECHO   ^> make TestGetPaymentOrdersByAuthorID                - Run TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID
	@ECHO   ^> make TestGetAllPaymentOrders                       - Run TestReferralRepositoryTestSuite/TestGetAllPaymentOrders
	@ECHO   ^> make TestCreatePaymentOrder                        - Run TestReferralRepositoryTestSuite/TestCreatePaymentOrder
	@ECHO   ^> make TestGetPaymentOrdersByAuthorID_Empty          - Run TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID_Empty
	@ECHO   ^> make TestAddTrHashToPaymentOrder                    - Run TestReferralRepositoryTestSuite/TestAddTrHashToPaymentOrder
	@ECHO   ^> make help                                          - Display this help information

referral-test:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite'

TestGetPaymentOrdersByAuthorID:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID'

TestGetAllPaymentOrders:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetAllPaymentOrders'

TestCreatePaymentOrder:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestCreatePaymentOrder'

TestGetPaymentOrdersByAuthorID_Empty:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetPaymentOrdersByAuthorID_Empty'

TestDeletePaymentOrder:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestDeletePaymentOrder'

TestGetDebtFromAuthorToReferrer:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestGetDebtFromAuthorToReferrer'

TestAddTrHashToPaymentOrder:
	go test -v ../test/referral/repository -run 'TestReferralRepositoryTestSuite/TestAddTrHashToPaymentOrder'