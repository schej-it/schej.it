import { serverURL, errors } from '@/constants'

/* 
  Fetch utils
*/
export const get = (route) => {
  return fetchMethod('GET', route)
}

export const post = (route, body = {}) => {
  return fetchMethod('POST', route, body)
}

export const patch = (route, body = {}) => {
  return fetchMethod('PATCH', route, body)
}

export const _delete = (route, body = {}) => {
  return fetchMethod('DELETE', route, body)
}

export const fetchMethod = (method, route, body = {}) => {
  /* Calls the given route with the give method and body */
  const params = {
    method,
    credentials: 'include',
  }

  if (method !== 'GET') {
    // Add params specific to POST/PATCH/DELETE
    params.headers = {
      'Content-Type': 'application/json',
    }
    params.body = JSON.stringify(body)
  }

  return fetch(serverURL + route, params)
    .then(async (res) => {
      const text = await res.text()

      // Check if response was ok
      if (!res.ok) {
        throw JSON.parse(text)
      }

      // Parse data if it is json, otherwise throw an error
      try {
        return JSON.parse(text)
      } catch (err) {
        throw { error: errors.JsonError }
      }
    })
    .then((data) => {
      return data
    })
}