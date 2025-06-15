import { post, _delete } from "./fetch_utils"

const FOLDER_API_ROUTE = "/user/folders"

export const createFolder = (name, color) => {
  return post(FOLDER_API_ROUTE, { name, color })
}

export const deleteFolder = (folderId) => {
  return _delete(`${FOLDER_API_ROUTE}/${folderId}`)
}

export const moveEventIntoFolder = (eventId, folderId) => {
  return post(`${FOLDER_API_ROUTE}/${folderId}/add-event`, { eventId })
}
