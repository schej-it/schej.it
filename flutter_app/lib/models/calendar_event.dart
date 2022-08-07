import 'dart:collection';

import 'package:table_calendar/table_calendar.dart';

// CalendarEvent contains data for a single event
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

// CalendarEvents stores CalendarEvents and allows you to access them by the
// day they occur on
class CalendarEvents {
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

  List<CalendarEvent> getEventsForDay(DateTime day) {
    final events = _eventsByDay[day];
    if (events == null) {
      return <CalendarEvent>[];
    }
    return events;
  }

  int _getHashCode(DateTime date) {
    return date.toLocal().toIso8601String().substring(0, 10).hashCode;
  }
}
