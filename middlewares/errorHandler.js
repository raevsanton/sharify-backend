import refreshTokenAndRetry from "../utils/refreshTokenAndRetry.js";

const errorHandler = async (error, request, response, next) => {
  if (error.status === 401) {
    await refreshTokenAndRetry(request, response);
  }
  next()
}

export default errorHandler;
