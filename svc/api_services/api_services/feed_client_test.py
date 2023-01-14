import os
from feed_client import FeedClient

if __name__ == '__main__':
    host = os.environ.get('FEED_MANAGER_SERVICE_HOST', 'localhost')
    port = int(os.environ.get('FEED_MANAGER_SERVICE_PORT', '6060'))
    cli =FeedClient(host, port)
    resp = cli.get_feed('ani')
    print(resp)