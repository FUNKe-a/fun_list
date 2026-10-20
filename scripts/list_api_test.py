import pytest
import requests

def test_list_directors(url_prefix, created_director):
    url = url_prefix + "/directors"
    response = requests.get(url)
    
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)
    assert len(data) >= 1
    assert any(d["DirectorID"] == created_director for d in data)

def test_list_media(url_prefix, created_media):
    url = url_prefix + "/media"
    response = requests.get(url)
    
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)
    assert len(data) >= 1
    assert any(m["MediaID"] == created_media for m in data)

def test_list_comments(url_prefix, created_comment):
    url = url_prefix + "/comments"
    response = requests.get(url)
    
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)
    assert len(data) >= 1
    assert any(c["CommentID"] == created_comment for c in data)

