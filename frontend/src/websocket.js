import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import store from './store'

Vue.use(VueNativeSock, `//${document.location.host}/ws`, {
  connectManually: true,
  store: store,
  passToStoreHandler: (eventName, event, nextHandler) => {
    if (eventName.toUpperCase() == 'SOCKET_ONMESSAGE' && event.data) {
      event.data.split( '\n' ).forEach( text => {
        if (text != "") {
          store.dispatch( 'APPEND_DATA', JSON.parse(text) )
            .catch( error => console.log("Error in ws handler: ", error) )
        }
      })
    } else {
      nextHandler( eventName, event )
    }
  }
})