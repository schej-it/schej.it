import { post, _delete } from "./fetch_utils"

const FOLDER_API_ROUTE = "/user/folders"

export const createFolder = (name, color) => {
  return post(FOLDER_API_ROUTE, { name, color })
}

export const deleteFolder = (folderId) => {
  return _delete(`${FOLDER_API_ROUTE}/${folderId}`)
}

export const setEventFolder = (eventId, folderId) => {
  return post(`/user/events/${eventId}/set-folder`, { folderId })
}
