import json
import os

import numpy as np
import pytest
from app.index import index
from app.main import _app as application
from app.tests.data import INPUT_VECTORS, TEST_VECTOR


@pytest.fixture
def app():
    # As part of the very first test, delete all local files
    try:
        os.remove(application.config['INDEX_FILE'])
        os.remove(application.config['INDEX_LOCK_FILE'])
        os.remove(application.config['INDEX_STATUS_FILE'])
    except FileNotFoundError:
        # Ignore if the files are already removed
        pass
    return application


@pytest.mark.run(order=0)
def test_health(client):
    """
    Test /health/ endpoint
    """
    response = client.get('/health/')
    assert response.status_code == 200
    assert response.json == {}


@pytest.mark.run(order=1)
def test_search_invalid(client):
    """
    Test invalid search requests
    """
    response = client.post('/search/')
    assert response.status_code == 400
    assert response.json == {'reason': 'Malformed request body'}

    response = client.post('/search/', data=json.dumps({}))
    assert response.status_code == 400
    assert response.json == {'reason': 'Malformed request body'}

    response = client.post('/search/', headers={'Content-Type': 'application/json'}, data=json.dumps({}))
    assert response.status_code == 400
    assert response.json == {'reason': '`vector` should be an array'}

    response = client.post('/search/', headers={'Content-Type': 'application/json'}, data=json.dumps({'vector': {}}))
    assert response.status_code == 400
    assert response.json == {'reason': '`vector` should be an array'}

    response = client.post('/search/', headers={'Content-Type': 'application/json'}, data=json.dumps({'vector': []}))
    assert response.status_code == 400
    assert response.json == {'reason': '`vector` should contain 512 entries'}


@pytest.mark.run(order=2)
def test_search_on_empty_index(client):
    """
    Test on an empty index
    """
    # Empty the index
    index().remove(np.arange(512))

    # Perform the search
    response = client.post('/search/',
                           headers={'Content-Type': 'application/json'},
                           data=json.dumps({'vector': TEST_VECTOR}))
    assert response.status_code == 404
    assert response.json == {'reason': 'No matching item is found', 'score': -1.0}


@pytest.mark.run(order=3)
def test_search(client):
    """
    Test a successful search with a properly registered index
    """
    # Populate the index with a `/register/` call
    response = client.post('/register/',
                           headers={'Content-Type': 'application/json'},
                           data=json.dumps({'vectors': [
                               [i, i, v] for i, v in enumerate(INPUT_VECTORS)
                           ]}))
    assert response.status_code == 200

    for _ in range(3):
        response = client.post('/search/',
                               headers={'Content-Type': 'application/json'},
                               data=json.dumps({'vector': TEST_VECTOR}))
        assert response.status_code == 200
        assert response.json == {'faiss_id': 36, 'score': 1309344.5, 'user_uid': 'NEW-FACE'}

    # Remove vector ID 36 from the index with a `/deregister/` call
    response = client.post('/deregister/',
                           headers={'Content-Type': 'application/json'},
                           data=json.dumps({'vector_ids': [36]}))
    assert response.status_code == 200
    assert response.json == {}

    # Perform the exact same search as above, the result is no longer 36
    for _ in range(3):
        response = client.post('/search/',
                               headers={'Content-Type': 'application/json'},
                               data=json.dumps({'vector': TEST_VECTOR}))
        assert response.status_code == 200
        assert response.json == {'faiss_id': 0, 'score': 1304253.0, 'user_uid': 'NEW-FACE'}
