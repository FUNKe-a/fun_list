import pytest
import requests

def test_delete_director(url_prefix, created_director):
    url = f"{url_prefix}/directors/{created_director}"
    response = requests.delete(url)
    assert response.status_code in [200, 204]

    assert requests.get(url).status_code == 404

def test_delete_comment(url_prefix, created_comment):
    url = f"{url_prefix}/comments/{created_comment}"
    response = requests.delete(url)
    assert response.status_code in [200, 204]

    assert requests.get(url).status_code == 404

def test_delete_media(url_prefix, created_media):
    url = f"{url_prefix}/media/{created_media}"
    response = requests.delete(url)
    assert response.status_code in [200, 204]
    
    assert requests.get(url).status_code == 404

