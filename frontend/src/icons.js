// icons.js

import Vue from 'vue'
import { library } from '@fortawesome/fontawesome-svg-core'

// import individual icons and add to library
import {
  faBan,
  faCoffee,
  faCheck,
  faExclamationTriangle,
  faInfo,
  faSync,
  faTachometerAlt,
  faEdit,
  faTrash,
  faDownload, faPause, faPlay, faStop
} from '@fortawesome/free-solid-svg-icons'

import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faCoffee)
library.add(faTachometerAlt)
library.add(faCheck)
library.add(faInfo)
library.add(faBan)
library.add(faSync)
library.add(faExclamationTriangle)
library.add(faEdit)
library.add(faTrash, faDownload, faPause, faPlay, faStop)

Vue.component('font-awesome-icon', FontAwesomeIcon)
