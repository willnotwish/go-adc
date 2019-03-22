import Auth0 from 'auth0-js'
import EventEmitter from 'events'
import Store from '../store'

const webAuth = new Auth0.WebAuth({
  domain:        process.env.VUE_APP_AUTH0_DOMAIN,
  redirectUri:  `${window.location.origin}/callback`,
  clientID:      process.env.VUE_APP_AUTH0_CLIENT_ID,
  audience:      process.env.VUE_APP_AUTH0_AUDIENCE,
  scope:        'openid profile email',
  responseType: 'token id_token',
})

const localStorageKey = 'loggedIn'
const loginEvent = 'loginEvent'

class AuthService extends EventEmitter {

  idToken = null

  tokenExpiry = null

  profile = null

  accessToken = null

  accessTokenExpiry = null

  // Starts the user login flow
  login(customState) {
    webAuth.authorize({
      appState: customState
    })
  }

  // Handles the callback request from Auth0
  handleAuthentication() {
    return new Promise((resolve, reject) => {
      webAuth.parseHash((err, authResult) => {
        if (err) {
          reject(err)
        } else {
          this.localLogin(authResult)
          resolve(authResult.idToken)
        }
      })
    })
  }

  localLogin(authResult) {
    this.idToken = authResult.idToken;
    this.profile = authResult.idTokenPayload;

    // Convert the JWT expiry time from seconds to milliseconds
    this.tokenExpiry = new Date(this.profile.exp * 1000)

    // Save the Access Token and expiry time in memory
    this.accessToken = authResult.accessToken;

    // Convert expiresIn to milliseconds and add the current time
    // (expiresIn is a relative timestamp, but an absolute time is desired)
    this.accessTokenExpiry = new Date(Date.now() + authResult.expiresIn * 1000);

    localStorage.setItem(localStorageKey, 'true')

    this.emit(loginEvent, {
      loggedIn: true,
      profile: authResult.idTokenPayload,
      state: authResult.appState || {}
    })
  }

  logOut() {
    localStorage.removeItem(localStorageKey)

    this.idToken = null
    this.tokenExpiry = null
    this.profile = null

    webAuth.logout({
      returnTo: window.location.origin
    })

    this.emit(loginEvent, { loggedIn: false })
  }

  isAuthenticated() {
    return Date.now() < this.tokenExpiry &&
      localStorage.getItem(localStorageKey) === 'true'
  }

  renewTokens() {
    return new Promise((resolve, reject) => {
      if (localStorage.getItem(localStorageKey) !== "true") {
        return reject("Cannot renew tokens because the user is not logged in")
      }

      webAuth.checkSession({}, (err, authResult) => {
        if (err) {
          reject(err)
        } else {
          this.localLogin(authResult)
          resolve(authResult)
        }
      })
    })
  }

  isAccessTokenValid() {
    return this.accessToken && this.accessTokenExpiry && Date.now() < this.accessTokenExpiry
  }

  getAccessToken() {
    return new Promise((resolve, reject) => {
      if (this.isAccessTokenValid()) {
        resolve(this.accessToken)
      } else {
        this.renewTokens().then(authResult => {
          resolve(authResult.accessToken)
        }, reject)
      }
    })
  }
}

export default new AuthService()