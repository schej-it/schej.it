// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyBZYjd-MK7OApHtiOIINZK2mCZOVJ3MpMI",
  authDomain: "schej-it.firebaseapp.com",
  projectId: "schej-it",
  storageBucket: "schej-it.appspot.com",
  messagingSenderId: "523323684219",
  appId: "1:523323684219:web:bfaae11213b8d1cfd4b9de",
  measurementId: "G-CRZP0XW10L"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const analytics = getAnalytics(app);