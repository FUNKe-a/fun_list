import pytest
import requests

def test_get_director_by_id(url_prefix, created_director):
    url = f"{url_prefix}/directors/{created_director}"
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

def test_get_media_by_id(url_prefix, created_media):
    url = f"{url_prefix}/media/{created_media}"
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

def test_get_comments_by_id(url_prefix, created_comment):
    url = f"{url_prefix}/comments/{created_comment}"
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
