import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/router/app_router.gr.dart';

class EventsPage extends StatelessWidget {
  const EventsPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(title: 'My events', isRoot: true),
      body: Center(
        child: Column(
          children: <Widget>[
            const Text('My events hehe'),
            ElevatedButton(
              onPressed: () {
                AutoRouter.of(context).push(EventPageRoute(eventId: '1234'));
              }, 
              child: const Text('Go to specific event!'),
            )
          ], 
        ),
      ), 
    ); 
  }
}

