// Event contains data for a single event
class Event {
  final String id;
  final String ownerId;
  final String name;
  final DateTime startDate;
  final DateTime endDate;
  final List<String> responses;

  const Event({
    required this.id,
    required this.ownerId,
    required this.name,
    required this.startDate,
    required this.endDate,
    required this.responses,
  });

  @override
  String toString() {
    return '{Event id:"$id" ownerId:"$ownerId" title:"$name" startDate:$startDate endDate:$endDate participants:$responses}';
  }
}

// Events stores Event objects that lets you access by created/joined and month
// class Events {
//   late final Map<String, List<Events>> _createdEventsByMonth;
//   late final Map<String, List<Events>> _joinedEventsByMonth;

//   Events({
//     required List<Event> events,
//   }) {
//     _createdEventsByMonth = LinkedHashMap(equals: isSameDay, hashCode: _getHashCode);

//     for (CalendarEvent event in events) {
//       if (_eventsByDay[event.startDate] == null) {
//         _eventsByDay[event.startDate] = [event];
//       } else {
//         _eventsByDay[event.startDate]!.add(event);
//       }
//     }
//   }

//   Map<DateTime, List<CalendarEvent>> get eventsByDay => _eventsByDay;

//   int _getHashCode(DateTime date) {
//     return date.toLocal().toIso8601String().substring(0, 10).hashCode;
//   }
// }
