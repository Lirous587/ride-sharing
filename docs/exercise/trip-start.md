Exercise

With all of this new information over the last couple of lessons, let's put everything to practice and create the full communication flow code for the /trip/start endpoint. 


 This is the endpoint the will be called once the user selects a ride fare or package for the desired route.

So the exercise is:

- Implement the HTTP signature and handler on the gateway service.

- Parse the payload. (I provide this)

- Add the needed protobuf definition on the trip.proto for the new method. 

- Implement the TripService interface and create an empty handler (leave it empty implementation for now). 

- Finally, call the gRPC CreateTrip method from the HTTP handler on the gateway service.



Note: Get the exercise starter code if you want and make sure to run the make generate-proto