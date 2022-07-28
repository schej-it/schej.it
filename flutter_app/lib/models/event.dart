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
