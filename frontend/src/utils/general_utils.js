/*
  General utils
*/

var timeoutId
/** Calls callback() on long press */
export const onLongPress = (element, callback, capture = false) => {

  element.addEventListener(
    'touchstart',
    function (e) {
      timeoutId = setTimeout(function () {
        timeoutId = null
        e.stopPropagation()
        callback(e.target)
      }, 500)
    },
    capture
  )

  element.addEventListener(
    'contextmenu',
    function (e) {
      e.preventDefault()
    },
    capture
  )

  element.addEventListener(
    'touchend',
    function () {
      if (timeoutId) clearTimeout(timeoutId)
    },
    capture
  )

  element.addEventListener(
    'touchmove',
    function () {
      if (timeoutId) clearTimeout(timeoutId)
    },
    capture
  )
}

/** Returns whether the given value is between lower and upper */
export const isBetween = (value, lower, upper, inclusive = true) => {
  if (inclusive) {
    return value >= lower && value <= upper
  } else {
    return value > lower && value < upper
  }
}

/** Clamps the given value between the given ranges */
export const clamp = (value, lower, upper) => {
  if (value < lower) return lower
  if (value > upper) return upper
  return value
}

export const isPhone = (vuetify) => {
  return vuetify.breakpoint.name === 'xs'
}

export const br = (vuetify, breakpoint) => {
  return vuetify.breakpoint.name === breakpoint
}

/** convert base64 to raw binary data held in a string */
export const dataURItoBlob = (dataURI) => {
  // doesn't handle URLEncoded DataURIs - see SO answer #6850276 for code that does this
  var byteString = atob(dataURI.split(',')[1])

  // separate out the mime component
  var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]

  // write the bytes of the string to an ArrayBuffer
  var ab = new ArrayBuffer(byteString.length)
  var ia = new Uint8Array(ab)
  for (var i = 0; i < byteString.length; i++) {
    ia[i] = byteString.charCodeAt(i)
  }

  return new Blob([ab], { type: mimeString })
}

/** Reformats the given event object to the format we want */
export const processEvent = (event) => {
  if (event.startDate) {
    event.startDate = new Date(event.startDate)
    event.endDate = new Date(event.endDate)
    event.startTime = event.startDate.getHours()
    event.endTime = event.endDate.getHours()
  } 
}