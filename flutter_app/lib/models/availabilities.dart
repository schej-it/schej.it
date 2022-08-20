import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/utils.dart';
import 'package:sorted_list/sorted_list.dart';

class Availabilities {
  // Calculates and returns a list of availability blocks for all the provided
  // user events for the given day
  static List<AvailabilityBlock> getUsersAvailabilityForDay(
    DateTime day,
    Map<String, List<CalendarEvent>> userEvents,
  ) {
    final initialAvailabilities = <AvailabilityBlock>[];
    for (final userId in userEvents.keys) {
      initialAvailabilities.addAll(parseDayAvailabilityFromCalendarEvents(
        day,
        userId,
        CalendarEvents(events: userEvents[userId]!),
      ));
    }
    return _calculateAvailabilities(initialAvailabilities);
  }

  // Calculates and returns an initial list of availability blocks from calendar
  // events for the given day
  static List<AvailabilityBlock> parseDayAvailabilityFromCalendarEvents(
    DateTime day,
    String userId,
    CalendarEvents calendarEvents,
  ) {
    final availability = <AvailabilityBlock>[];

    // If there are no events, then user is available all day
    final events = calendarEvents.events;
    if (events.isEmpty) {
      availability.add(AvailabilityBlock(
        startDate: getDateWithTime(day, 0, local: true),
        endDate: getDateWithTime(day, 23.99, local: true),
        usersAvailable: {userId},
      ));
      return availability;
    }

    DateTime curDateTime = getDateWithTime(day, 0, local: true);
    for (final event in events) {
      if (event.startDate.isBefore(curDateTime)) {
        curDateTime = event.endDate;
        continue;
      }
      final block = AvailabilityBlock(
        startDate: curDateTime,
        endDate: event.startDate,
        usersAvailable: {userId},
      );
      if (!_isEmptyAvailabilityBlock(block)) availability.add(block);
      curDateTime = event.endDate;
    }

    final block = AvailabilityBlock(
      startDate: curDateTime,
      endDate: getDateWithTime(day, 23.99, local: true),
      usersAvailable: {userId},
    );
    if (!_isEmptyAvailabilityBlock(block)) availability.add(block);

    return availability;
  }

  // Calculates and returns overlapping availabilities from initial availabilities
  // parsed from calendar events
  static List<AvailabilityBlock> _calculateAvailabilities(
    List<AvailabilityBlock> initialAvailabilities,
  ) {
    final availabilities = <AvailabilityBlock>[];
    final arr = SortedList<AvailabilityBlock>(_sortByStartDate);
    arr.addAll(initialAvailabilities);

    // print('arr length: ${arr.length}');
    // int i = 0;
    while (arr.isNotEmpty) {
      // i++;

      final a = arr[0];
      if (arr.length == 1) {
        availabilities.add(a);
        arr.removeAt(0);
        continue;
      }

      final b = arr[1];

      if (!_doesOverlap(a, b)) {
        availabilities.add(a);
        arr.removeAt(0);
        continue;
      } else {
        final overlap = _getOverlap(a, b);

        // Remove a and b
        arr.removeAt(0);
        arr.removeAt(0);

        arr.add(overlap);

        final diffA = _getDiff(a, overlap);
        final diffB = _getDiff(b, overlap);
        arr.addAll(diffA);
        arr.addAll(diffB);
      }
    }
    // print('iterations: $i');

    return availabilities;
  }

  static int _sortByStartDate(a, b) => a.startDate.compareTo(b.startDate);

  // Returns whether a and b overlap
  static bool _doesOverlap(AvailabilityBlock a, AvailabilityBlock b) {
    return a.startDate.difference(b.startDate).inMinutes <= 0 &&
        a.endDate.difference(b.startDate).inMinutes > 0;
  }

  // Returns the overlap between a and b
  // Assumptions:
  // - a starts before b
  // - a and b indeed overlap
  static AvailabilityBlock _getOverlap(
      AvailabilityBlock a, AvailabilityBlock b) {
    // Merge userAvailable array
    final usersAvailable = Set<String>.from(a.usersAvailable)
      ..addAll(b.usersAvailable);

    return AvailabilityBlock(
      startDate: b.startDate,
      endDate: a.endDate.isBefore(b.endDate) ? a.endDate : b.endDate,
      usersAvailable: usersAvailable,
    );
  }

  // Returns the difference between a and b
  static List<AvailabilityBlock> _getDiff(
      AvailabilityBlock a, AvailabilityBlock b) {
    final diff = <AvailabilityBlock>[];

    final first = AvailabilityBlock(
      startDate: a.startDate,
      endDate: b.startDate,
      usersAvailable: a.usersAvailable,
    );
    if (!_isEmptyAvailabilityBlock(first)) diff.add(first);

    final second = AvailabilityBlock(
      startDate: b.endDate,
      endDate: a.endDate,
      usersAvailable: a.usersAvailable,
    );
    if (!_isEmptyAvailabilityBlock(second)) diff.add(second);

    return diff;
  }

  // Returns whether the availability is empty, i.e. the time range is either
  // empty or invalid
  static bool _isEmptyAvailabilityBlock(AvailabilityBlock a) {
    return a.startDate.difference(a.endDate).inMinutes >= 0;
  }

  /* 
  PSEUDOCODE: 
  let arr = availability block array that is kept sorted

  while there are still availability blocks
    let i = first availability block
    let j = second availability block
    if i and j dont overlap, move i to availabilities array and return 
    if i and j do overlap 
      create a new availability block from the overlap and merge the usersAvailable arrays, insert it into array
      change i from availability blocks array to be i - overlap (not necessarily change end time of i to be start time of overlap because what if j is fully contained by i) 
      change j to be j - overlap
      reinsert i and j so sorting updates
      if either i or j is a no time event now (i.e. startDate == endDate), remove it from the blocks
  
  */
}

class AvailabilityBlock {
  final DateTime startDate;
  final DateTime endDate;
  final Set<String> usersAvailable;

  double get startTime => startDate.hour + startDate.minute / 60;
  double get endTime => endDate.hour + endDate.minute / 60;

  AvailabilityBlock({
    required this.startDate,
    required this.endDate,
    required this.usersAvailable,
  });

  @override
  String toString() {
    return '{AvailabilityBlock startDate:$startDate endDate:$endDate usersAvailable:$usersAvailable}';
  }
}
