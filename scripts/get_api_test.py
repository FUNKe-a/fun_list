import pytest
import requests

url_prefix = "http://127.0.0.1:8080/api";

def test_get_director_by_id():
    url = url_prefix + "/directors/1"
    response = requests.get(url)

    assert response.status_code == 200

    assert response.headers["Content-Type"] == "application/json"

    data = response.json()  
    assert isinstance(data, dict)
    assert "DirectorID" in data
    assert "Name" in data
    assert "DateOfBirth" in data
    assert "DateOfDeath" in data
    assert "Biography" in data

def test_get_media_by_id():
    url = url_prefix + "/media/1"
    response = requests.get(url)

    assert response.status_code == 200

    assert response.headers["Content-Type"] == "application/json"

    data = response.json()  
    assert isinstance(data, dict)
    assert "MediaID" in data
    assert "Title" in data
    assert "DirectorID" in data
    assert "ReleasedAt" in data
    assert "Synopsis" in data

def test_get_comments_by_id():
    url = url_prefix + "/comments/1"
    response = requests.get(url)

    assert response.status_code == 200

    assert response.headers["Content-Type"] == "application/json"

    data = response.json()  
    assert isinstance(data, dict)
    assert "CommentID" in data
    assert "UserID" in data
    assert "MediaID" in data
    assert "Content" in data
    assert "Rating" in data
    assert "CreatedAt" in data
