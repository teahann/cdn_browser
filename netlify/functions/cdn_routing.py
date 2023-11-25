import json
import http.client
from urllib.parse import urlparse

def handler(event, context):
  path = event.get('path')
  
  if path:
    base_url = 'https://dchqsxvlazultvsphmau.supabase.co/storage/v1/object/public/'
    imgUrl = f"{base_url}{path}"
    parsed_url = urlparse(imgUrl)
    connection = http.client.HTTPSConnection(parsed_url.netloc)
    connection.request('GET', parsed_url.path)
