import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';

class EventPage extends StatelessWidget {
  const EventPage({Key? key, @PathParam('id') required this.eventId}) : super(key: key);

  final String eventId;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Event'),
        leading: const AutoLeadingButton(),
      ),
      body: Center(
        child: Text('Specific event, event id: $eventId'),
      )
    ); 
   }
}

