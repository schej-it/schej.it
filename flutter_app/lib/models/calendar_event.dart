import 'dart:collection';

import 'package:table_calendar/table_calendar.dart';

class CalendarEvent {
  final String title;
  final DateTime startDate;
  final DateTime endDate;

  double get startTime => startDate.hour + startDate.minute / 60;
  double get endTime => endDate.hour + endDate.minute / 60;

  const CalendarEvent({
    required this.title,
    required this.startDate,
    required this.endDate,
  });

  @override
  String toString() {
    return '{CalendarEvent title:"$title" startDate:$startDate endDate:$endDate}';
  }
}

class CalendarEvents {
  // A map mapping a string representing the day ("7-15-2022") to the events
  // occurring on that day
  late final Map<DateTime, List<CalendarEvent>> _eventsByDay;

  CalendarEvents({
    required List<CalendarEvent> events,
  }) {
    _eventsByDay = LinkedHashMap(equals: isSameDay, hashCode: _getHashCode);

    for (CalendarEvent event in events) {
      if (_eventsByDay[event.startDate] == null) {
        _eventsByDay[event.startDate] = [event];
      } else {
        _eventsByDay[event.startDate]!.add(event);
      }
    }
  }

  Map<DateTime, List<CalendarEvent>> get eventsByDay => _eventsByDay;

  int _getHashCode(DateTime date) {
    return date.toLocal().toIso8601String().substring(0, 10).hashCode;
  }
}
