import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';

class EventPage extends StatelessWidget {
  const EventPage({Key? key, @PathParam('id') required this.eventId}) : super(key: key);

  final String eventId;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(titleString: 'Event'),
      body: Center(
        child: Text('Specific event, event id: $eventId'),
      )
    ); 
   }
}

