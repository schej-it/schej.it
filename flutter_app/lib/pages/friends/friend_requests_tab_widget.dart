import 'package:flutter/material.dart';
import 'package:flutter_app/components/friend_request_card.dart';
import 'package:flutter_app/constants/constants.dart';

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
    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: SchejConstants.pagePadding.left,
      ),
      child: _buildFriendRequestCards(),
    );
  }

  Widget _buildFriendRequestCards() {
    return ListView.builder(
      itemCount: friendRequests.length,
      itemBuilder: (context, index) {
        final request = friendRequests[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: FriendRequestCard(
            name: request['name'] as String,
            requestTimestamp: request['requestTimestamp'] as DateTime,
          ),
        );
      },
    );
  }
}
