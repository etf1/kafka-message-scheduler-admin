const HEADER_APP_JSON = {
  "Content-Type": "application/json",
  Accept: "application/json",
};

const SERVER_ERR_MSG = "Bad response from server";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const post = async <T = any>(
  url: string,
  body: string,
  header?: HeadersInit
): Promise<T> => {
  const response = await fetch(url, {
    method: "POST",
    body: body,
    headers: header || HEADER_APP_JSON,
  });

  if (response.status >= 400) {
    throw new Error(SERVER_ERR_MSG);
  }
  return response.json();
};
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const get = async <T = any>(
  url: string,
  header?: HeadersInit
): Promise<T> => {
  const response = await fetch(url, {
    method: "GET",
    headers: header || HEADER_APP_JSON,
  });
  if (response.status >= 400) {
    throw new Error(SERVER_ERR_MSG);
  }
  return response.json();
};
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const del = async <T = any>(
  url: string,
  header?: HeadersInit
): Promise<T> => {
  const response = await fetch(url, {
    method: "DELETE",
    headers: header || HEADER_APP_JSON,
  });
  if (response.status >= 400) {
    throw new Error(SERVER_ERR_MSG);
  }
  return response.json();
};
