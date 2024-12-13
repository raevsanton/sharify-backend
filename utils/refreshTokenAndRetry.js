import axios from 'axios';

const refreshTokenAndRetry = async (request, response) => {
  try {
    const refreshToken = request.headers.authorization.split(' ').at(2);

    const data = await axios.post(
      'https://accounts.spotify.com/api/token',
      new URLSearchParams({
        grant_type: 'refresh_token',
        refresh_token: refreshToken,
      }),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Authorization': `Basic ${(new Buffer.from(CLIENT_ID + ':' + CLIENT_SECRET).toString('base64'))}`
        }
      }
    );

    const newAccessToken = data.data.access_token;
    const { protocol, originalUrl, method, body } = request;

    const newResponse = await axios({
      url: `${protocol}://${request.get('host')}${originalUrl}`,
      method,
      headers: {
          ...request.get('headers'),
          Authorization: `Bearer ${newAccessToken}`,
          'Content-Type': 'application/json',
      },
        data: body,
    });

    const newResponseData = await newResponse.json();
    return response.send(newResponseData);
  } catch (error) {
      return response.send(error);
  }
};

export default refreshTokenAndRetry;
