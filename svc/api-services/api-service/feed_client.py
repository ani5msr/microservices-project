import grpc
import feed_pb2
import feed_pb2_grpc


class FeedClient(object):
    """
    Client for accessing the gRPC functionality
    """

    def __init__(self, host, port):
        # instantiate a communication channel
        self.channel = grpc.insecure_channel(
            '{}:{}'.format(host, port))

        # bind the client to the server channel
        self.stub = feed_pb2_grpc.FeedStub(self.channel)

    def get_feed(self, username, startToken=None):
        """
        Client function to call the rpc for GetDigest
        """
        req = feed_pb2.GetFeedRequest(username=username, startToken=startToken)
        return self.stub.GetFeed(req)
