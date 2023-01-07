/** Redirects user to the correct google sign in page */
export const signInGoogle = (state = null, consent = false) => {
  const clientId = '523323684219-jfakov2bgsleeb6den4ktpohq4lcnae2.apps.googleusercontent.com'
  const redirectUri = `${window.location.origin}/auth`
  const scope = encodeURIComponent(
    'openid email profile https://www.googleapis.com/auth/calendar.calendarlist.readonly https://www.googleapis.com/auth/calendar.events.readonly'
  )

  let stateString = ''
  if (state !== null) {
    state = encodeURIComponent(JSON.stringify(state))
    stateString = `&state=${state}`
  }

  let promptString = ''
  if (consent) {
    promptString = '&prompt=consent'
  } else {
    promptString = '&prompt=select_account'
  }

  window.location.href = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${redirectUri}&response_type=code&scope=${scope}&access_type=offline${promptString}${stateString}`
}
