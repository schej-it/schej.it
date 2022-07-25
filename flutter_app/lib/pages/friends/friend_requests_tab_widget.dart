import 'package:flutter/material.dart';
import 'package:flutter/src/foundation/key.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter_app/components/friend_request_card.dart';

class FriendRequestsTabWidget extends StatefulWidget {
  const FriendRequestsTabWidget({Key? key}) : super(key: key);

  @override
  State<FriendRequestsTabWidget> createState() =>
      _FriendRequestsTabWidgetState();
}

class _FriendRequestsTabWidgetState extends State<FriendRequestsTabWidget> {
  // Variables
  var friendRequests = [
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {
      'name': 'Samantha Jones',
      'requestTimestamp': DateTime.now().subtract(const Duration(days: 2))
    },
    {
      'name': 'Tyler Smithson',
      'requestTimestamp': DateTime.now().subtract(const Duration(days: 4))
    },
  ];

  @override
  Widget build(BuildContext context) {
    return _buildFriendRequestCards();
  }

  Widget _buildFriendRequestCards() {
    final cards = <Widget>[];
    for (final request in friendRequests) {
      cards.add(Padding(
        padding: const EdgeInsets.only(bottom: 8),
        child: FriendRequestCard(
          name: request['name'] as String,
          requestTimestamp: request['requestTimestamp'] as DateTime,
        ),
      ));
    }

    return Column(children: cards);
  }
}
