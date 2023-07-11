/* utils for getting user's location */

export const getLocation = () => {
  return fetch("https://geolocation-db.com/json/", {
    method: "GET",
  }).then((res) => res.json())
}
