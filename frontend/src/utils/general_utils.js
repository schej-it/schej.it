/*
  General utils
*/

import { eventTypes } from "@/constants"
import { dateToDowDate, dateToTimeNum } from "./date_utils"

var timeoutId
/** Calls callback() on long press */
export const onLongPress = (element, callback, capture = false) => {
  element.addEventListener(
    "touchstart",
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
    "contextmenu",
    function (e) {
      e.preventDefault()
    },
    capture
  )

  element.addEventListener(
    "touchend",
    function () {
      if (timeoutId) clearTimeout(timeoutId)
    },
    capture
  )

  element.addEventListener(
    "touchmove",
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
  return vuetify.breakpoint.name === "xs"
}

export const br = (vuetify, breakpoint) => {
  return vuetify.breakpoint.name === breakpoint
}

/** convert base64 to raw binary data held in a string */
export const dataURItoBlob = (dataURI) => {
  // doesn't handle URLEncoded DataURIs - see SO answer #6850276 for code that does this
  var byteString = atob(dataURI.split(",")[1])

  // separate out the mime component
  var mimeString = dataURI.split(",")[0].split(":")[1].split(";")[0]

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
  let startDate = event.dates[0]
  if (event.type === eventTypes.DOW) {
    startDate = dateToDowDate(event.dates, startDate, 0, true)
  }

  event.startTime = dateToTimeNum(new Date(startDate), true)
  event.endTime = (event.startTime + event.duration) % 24
}

/** Checks whether email is a valid email */
export const validateEmail = (email) => {
  return String(email)
    .toLowerCase()
    .match(
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    )
}

/** Generates a group enabled calendar payload */
export const generateEnabledCalendarsPayload = (calendarAccounts) => {
  const payload = {}

  payload.guest = false
  payload.useCalendarAvailability = true
  payload.enabledCalendars = {}

  /** Determine which sub calendars are enabled */
  for (const email in calendarAccounts) {
    if (calendarAccounts[email].enabled) {
      payload.enabledCalendars[email] = []
      for (const subCalendarId in calendarAccounts[email].subCalendars) {
        if (calendarAccounts[email].subCalendars[subCalendarId].enabled) {
          payload.enabledCalendars[email].push(subCalendarId)
        }
      }
    }
  }

  return payload
}

/** Returns whether touch is enabled on the device */
export const isTouchEnabled = () => {
  return (
    "ontouchstart" in window ||
    navigator.maxTouchPoints > 0 ||
    navigator.msMaxTouchPoints > 0
  )
}

/** Returns whether the element is in the viewport */
export const isElementInViewport = (
  el,
  { topOffset = 0, leftOffset = 0, rightOffset = 0, bottomOffset = 0 }
) => {
  var rect = el.getBoundingClientRect()

  return (
    rect.top >= topOffset &&
    rect.left >= leftOffset &&
    rect.bottom <=
      (window.innerHeight || document.documentElement.clientHeight) +
        bottomOffset &&
    rect.right <=
      (window.innerWidth || document.documentElement.clientWidth) + rightOffset
  )
}
