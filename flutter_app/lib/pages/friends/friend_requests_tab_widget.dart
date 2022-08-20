import 'package:flutter/material.dart';
import 'package:flutter_app/components/friends/friend_request_card.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/friend_request.dart';
import 'package:provider/provider.dart';

class FriendRequestsTabWidget extends StatefulWidget {
  const FriendRequestsTabWidget({Key? key}) : super(key: key);

  @override
  State<FriendRequestsTabWidget> createState() =>
      _FriendRequestsTabWidgetState();
}

class _FriendRequestsTabWidgetState extends State<FriendRequestsTabWidget> {
  // Dummy test data.
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
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
    {'name': 'Winston Tilton', 'requestTimestamp': DateTime.now()},
  ];

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
        padding: const EdgeInsets.only(top: 10),
        child: Consumer<ApiService>(
          builder: (context, api, child) => _buildFriendRequestCards(api),
        ));
  }

  Widget _buildFriendRequestCards(ApiService api) {
    // CHECK: might not rerender when changed.
    List<FriendRequest> incomingFriendRequests =
        api.getIncomingFriendRequests();
    return ListView.builder(
      itemCount: incomingFriendRequests.length + 1,
      itemBuilder: (context, index) {
        if (index == incomingFriendRequests.length) {
          // Return sized box so FAB doesn't overlap
          return const SizedBox(height: 70);
        }

        final request = incomingFriendRequests[index];
        return Padding(
          padding: SchejConstants.pagePadding
              .copyWith(top: index == 0 ? 10 : 0, bottom: 10),
          child: FriendRequestCard(
            id: request.id,
            name: request.fromUser.fullName,
            requestTimestamp: request.createdAt,
            picture: request.fromUser.picture,
            api: context.read<ApiService>(),
          ),
        );
      },
    );
  }
}
