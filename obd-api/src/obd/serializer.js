import {Serializer as JSONApiSerializer} from 'jsonapi-serializer'

// Knows how to serialize from a pid value model

export default new JSONApiSerializer( 'pid-value', {
  attributes: ['pid', 'value'],
})