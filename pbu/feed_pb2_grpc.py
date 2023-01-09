# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import feed_pb2 as feed__pb2


class FeedStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetFeed = channel.unary_unary(
                '/pb.Feed/GetFeed',
                request_serializer=feed__pb2.GetFeedRequest.SerializeToString,
                response_deserializer=feed__pb2.GetFeedResponse.FromString,
                )


class FeedServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetFeed(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_FeedServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetFeed': grpc.unary_unary_rpc_method_handler(
                    servicer.GetFeed,
                    request_deserializer=feed__pb2.GetFeedRequest.FromString,
                    response_serializer=feed__pb2.GetFeedResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'pb.Feed', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Feed(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetFeed(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pb.Feed/GetFeed',
            feed__pb2.GetFeedRequest.SerializeToString,
            feed__pb2.GetFeedResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)