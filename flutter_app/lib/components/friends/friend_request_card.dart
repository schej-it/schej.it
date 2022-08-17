import 'package:flutter/material.dart';
import 'package:flutter_app/components/user_avatar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:intl/intl.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class FriendRequestCard extends StatelessWidget {
  final String id;
  final String name;
  final String picture;
  final DateTime requestTimestamp;
  final ApiService api;

  const FriendRequestCard(
      {Key? key,
      required this.id,
      required this.name,
      this.picture = 'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
      required this.requestTimestamp,
      required this.api})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: SchejConstants.listTileDecoration,
      child: ListTile(
        dense: true,
        leading: UserAvatar(src: picture),
        title: Text(name, style: SchejFonts.subtitle),
        subtitle: _buildRequestTimestamp(),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            IconButton(
              icon: const Icon(MdiIcons.close),
              color: SchejColors.black,
              onPressed: () {
                api.rejectFriendRequest(id);
              },
            ),
            IconButton(
              icon: const Icon(MdiIcons.check),
              color: SchejColors.darkGreen,
              onPressed: () {
                api.acceptFriendRequest(id);
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRequestTimestamp() {
    final dateString = DateFormat.MMMMd().format(requestTimestamp);
    final timestampString = 'Requested on $dateString';

    return Text(timestampString,
        style: SchejFonts.body.copyWith(color: SchejColors.pureBlack));
  }
}
