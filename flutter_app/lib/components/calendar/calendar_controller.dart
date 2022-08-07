import 'package:flutter/material.dart';
import 'package:flutter_app/models/calendar_event.dart';

// NOTE: UNUSED for now!
class CalendarController extends ChangeNotifier {
  final Map<String, CalendarEvents>? initialCalendarEvents;
  final DateTime? initialSelectedDay;
  final int? initialDaysVisible;
  final bool? initialShowEventTitles;
  final bool? initialShowAvatars;
  final String? initialActiveUser;

  late final Map<String, CalendarEvents> _calendarEvents;
  DateTime _selectedDay = DateTime.now();
  int _daysVisible = 3;
  bool _showEventTitles = true;
  bool _showAvatars = false;
  String? _activeUser;

  CalendarController({
    this.initialCalendarEvents,
    this.initialSelectedDay,
    this.initialDaysVisible,
    this.initialShowEventTitles,
    this.initialShowAvatars,
    this.initialActiveUser,
  }) {
    if (initialCalendarEvents == null) {
      _calendarEvents = <String, CalendarEvents>{};
    } else {
      _calendarEvents = initialCalendarEvents!;
    }

    if (initialSelectedDay != null) {
      _selectedDay = initialSelectedDay!;
    }

    if (initialDaysVisible != null) {
      _daysVisible = initialDaysVisible!;
    }

    if (initialShowEventTitles != null) {
      _showEventTitles = initialShowEventTitles!;
    }

    if (initialShowAvatars != null) {
      _showAvatars = initialShowAvatars!;
    }

    _activeUser = initialActiveUser;
  }

  DateTime get selectedDay => _selectedDay;
  set selectedDay(DateTime value) {
    _selectedDay = value;
    notifyListeners();
  }

  int get daysVisible => _daysVisible;
  set daysVisible(int value) {
    _daysVisible = value;
    notifyListeners();
  }

  bool get showEventTitles => _showEventTitles;
  set showEventTitles(bool value) {
    _showEventTitles = value;
    notifyListeners();
  }

  bool get showAvatars => _showAvatars;
  set showAvatars(bool value) {
    _showAvatars = value;
    notifyListeners();
  }

  String? get activeUser => _activeUser;
  set activeUser(String? value) {
    _activeUser = value;
    notifyListeners();
  }

  CalendarEvents? getCalendarEvents(String userId) {
    return _calendarEvents[userId];
  }

  void addCalendarEvents(String userId, CalendarEvents calendar) {
    _calendarEvents[userId] = calendar;
    notifyListeners();
  }

  void removeCalendarEvents(String userId) {
    _calendarEvents.remove(userId);
    notifyListeners();
  }
}
