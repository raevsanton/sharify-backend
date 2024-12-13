import express from 'express';
import axios from 'axios';
import cors from 'cors';
import errorHandler from './middlewares/errorHandler.js';
import 'dotenv/config';

var SPOTIFY_CLIENT_ID = process.env.SPOTIFY_CLIENT_ID
var CLIENT_SECRET = process.env.CLIENT_SECRET
var CLIENT_URL = process.env.CLIENT_URL

const app = express();

app.use(cors());
app.use(express.json())

const getUserId = async (token) => {
  const { data } = await axios.get('https://api.spotify.com/v1/me', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  return data.id;
}

const getLikedTracks = async (request, offset) => {
  const accessToken = request.headers.authorization.split(' ')[1];

  try {
    const { data } = await axios.get(`https://api.spotify.com/v1/me/tracks?offset=${offset}&limit=50`, {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
    });

    if (data.error) {
      const error = new Error(data.error?.message);
      error.status = data.error?.status;
      throw error;
    }

    const uris = data.items.map(({ track }) => track.uri);
    const total = data.total;

    return { total, uris };
  } catch (error) {
    throw error;
  }
};

const createPlaylist = async (token, { name, description, is_public }) => {
  const userID = await getUserId(token);

  try {
    const playlistData = await axios.post(
      `https://api.spotify.com/v1/users/${encodeURIComponent(userID)}/playlists`,
      {
        "name": name.length ? name.toString() : 'My Liked Songs',
        "description": description.toString(),
        "public": is_public,
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }
    );
    return playlistData.data.id;
  } catch (error) {
    next(error);
  }
}

const addTracksToPlaylist = async (token, uris, playlistID, position) => {
  try {
    const { data } = await axios.post(
      `https://api.spotify.com/v1/playlists/${playlistID}/tracks`,
      {
        "uris": uris,
        "position": position
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }
    );
    return data;
  } catch (error) {
    next(error);
  }
}

app.get('/auth', async (request, response) => {
  const code = request.query['code'];

  try {
    const { data } = await axios.post(
      'https://accounts.spotify.com/api/token',
      new URLSearchParams({
        grant_type: 'authorization_code',
        redirect_uri: CLIENT_URL,
        code,
      }),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        auth: {
          username: SPOTIFY_CLIENT_ID,
          password: CLIENT_SECRET,
        },
      }
    );

    const accessToken = data.access_token
    const refreshToken = data.refresh_token;

    response.send({
      access_token: accessToken,
      refresh_token: refreshToken,
    });
  } catch (error) {
    response.status(500).json(error);
  }
});

app.post('/playlist', async (request, response, next) => {
  try {
    const accessToken = request.headers.authorization.split(' ')[1];
    const { id: existPlaylistID, name, description, is_public } = request.body;

    let totalUris = 0;
    let accumulatedUris = [];
    let offset = 0;
    let position = 0;

    do {
      const { uris, total } = await getLikedTracks(request, offset);
      if (totalUris === 0) {
        totalUris = total;
      }

      offset += 50;
      accumulatedUris.push(...uris);
    } while (totalUris > accumulatedUris.length);

    const playlistID = existPlaylistID?.length
      ? existPlaylistID
      : await createPlaylist(accessToken, { name, description, is_public });

    for (let i = 0; i < accumulatedUris.length; i += 100) {
      const hundredUris = accumulatedUris.slice(i, i + 100);
      await addTracksToPlaylist(accessToken, hundredUris, playlistID, position);
      position += 100;
    }

    return response.json({
      playlist_id: playlistID,
      access_token: accessToken });
  } catch (error) {
    next(error);
  }
});

app.use(errorHandler)

app.listen(8080);
