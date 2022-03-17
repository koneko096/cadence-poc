package com.github.koneko096.cadencepoc.payment;

import io.grpc.stub.StreamObserver;
import org.lognet.springboot.grpc.GRpcService;

@GRpcService
public class Service extends PaymentGrpc.PaymentImplBase {
  @Override
  public void deductFare(Billing request, StreamObserver<Receipt> responseObserver) {
    Receipt res = Receipt.newBuilder()
        .setReceiptID(256)
        .build();
    responseObserver.onNext(res);
    responseObserver.onCompleted();
  }
}
