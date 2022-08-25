import 'dart:collection';

import 'package:flutter_app/utils.dart';
import 'package:json_annotation/json_annotation.dart';
import 'package:sorted_list/sorted_list.dart';
import 'package:table_calendar/table_calendar.dart';

part 'calendar_event.g.dart';

// CalendarEvent contains data for a single event
@JsonSerializable()
class CalendarEvent {
  @JsonKey(name: 'summary')
  final String title;
  DateTime startDate;
  DateTime endDate;

  double get startTime => startDate.hour + startDate.minute / 60;
  double get endTime => endDate.hour + endDate.minute / 60;

  CalendarEvent({
    required this.title,
    required this.startDate,
    required this.endDate,
  }) {
    startDate = startDate.toLocal();
    endDate = endDate.toLocal();
  }

  @override
  String toString() {
    return '{CalendarEvent title:"$title" startDate:$startDate endDate:$endDate}';
  }

  factory CalendarEvent.fromJson(Map<String, dynamic> json) =>
      _$CalendarEventFromJson(json);
  Map<String, dynamic> toJson() => _$CalendarEventToJson(this);
}

// CalendarEvents stores CalendarEvents and allows you to access them by the
// day they occur on
class CalendarEvents {
  late final Map<DateTime, List<CalendarEvent>> _eventsByDay;

  CalendarEvents({
    required List<CalendarEvent> events,
  }) {
    _eventsByDay = LinkedHashMap(equals: isSameDay, hashCode: _getHashCode);

    addEvents(events);
  }

  Iterable<DateTime> get days => _eventsByDay.keys;
  List<CalendarEvent> get events =>
      _eventsByDay.values.expand((i) => i).toList();

  List<CalendarEvent> getEventsForDay(DateTime day) {
    final events = _eventsByDay[day];
    if (events == null) {
      return <CalendarEvent>[];
    }
    return events;
  }

  void addEvents(List<CalendarEvent> newEvents) {
    // TODO: This won't necessarily remove duplicate events spanning multiple
    //       days since those events have been split and therefore won't be
    //       equal.
    // Remove duplicates
    final uniqueEventsSet = newEvents.toSet();
    final difference = uniqueEventsSet.difference(events.toSet());
    final uniqueEvents = difference.toList();

    while (uniqueEvents.isNotEmpty) {
      CalendarEvent event = uniqueEvents[0];
      uniqueEvents.removeAt(0);

      // Split events spanning multiple days into multiple events
      if (!isSameDay(event.startDate, event.endDate)) {
        DateTime curDate = event.startDate;
        while (!isSameDay(curDate, event.endDate)) {
          final splitEvent = CalendarEvent(
            title: event.title,
            startDate: curDate,
            endDate: getLocalDateWithTime(curDate, 23.99),
          );
          curDate =
              getLocalDateWithTime(curDate.add(const Duration(days: 1)), 0);
          uniqueEvents.add(splitEvent);
        }
        uniqueEvents.add(CalendarEvent(
          title: event.title,
          startDate: curDate,
          endDate: event.endDate,
        ));

        // Go to next iteration because event has been split up, so we shouldn't
        // add it
        continue;
      }

      // Add event to map
      if (_eventsByDay[event.startDate] == null) {
        _eventsByDay[event.startDate] =
            SortedList<CalendarEvent>(_sortByStartDate);
      }

      _eventsByDay[event.startDate.toLocal()]!.add(event);
    }
  }

  int _getHashCode(DateTime date) {
    return date.toLocal().toIso8601String().substring(0, 10).hashCode;
  }

  int _sortByStartDate(a, b) => a.startDate.compareTo(b.startDate);

  @override
  String toString() {
    return '{CalendarEvents eventsByDay:${_eventsByDay.toString()}}';
  }
}

// DayRange represents a day range
class DayRange {
  DateTime start;
  DateTime end;

  DayRange({
    required this.start,
    required this.end,
  });

  bool isInRange(DateTime day) {
    return compareToDay(day) == 0;
  }

  int compareToDay(DateTime day) {
    if (inRange(day, start, end)) {
      return 0;
    } else {
      return start.compareTo(day);
    }
  }

  @override
  String toString() {
    return '{DayRange start:$start end:$end }';
  }
}
